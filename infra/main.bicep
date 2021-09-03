param tags object = {
  evnironment: 'dev'
  costcode: '1234567890'
}
param adminGroupId string = 'f6a900e2-df11-43e7-ba3e-22be99d3cede'
param sshPublicKey string
param objectId string = '57963f10-818b-406d-a2f6-6e758d86e259'

var acrPullRoleId = resourceId('Microsoft.Authorization/roleDefinitions', '7f951dda-4ed3-4680-a7ca-43fe172d538d')
var suffix = '${uniqueString(resourceGroup().id)}'
var keyVaultName = 'kv-${suffix}'
var serviceBusName = 'sb-${suffix}'
var cosmosDbName = 'cs-${suffix}'
var acrName = 'acr${suffix}'
var vnetName = 'vnet-${suffix}'
var wksName = 'wks-${suffix}'
var aksName = 'aks-${suffix}'
var accessPolicies = [
  {
    permissions: {
      certificates: [
        'all'
      ]
      keys: [
        'all'
      ]
      secrets: [
        'all'
      ]
      storage: [
        'all'
      ]
    }
    objectId: objectId
    tenantId: subscription().tenantId
  }
]

module wksMod 'modules/wks.bicep' = {
  name: wksName
  params: {
    name: wksName
    retentionInDays: 30
    tags: tags
  }
}

module acrMod 'modules/acr.bicep' = {
  name: acrName
  params: {
    acrName: acrName
    tags: tags
  }
}

module kvMod 'modules/keyvault.bicep' = {
  name: 'keyVaultDeployment'
  params: {
    name: keyVaultName
    accessPolicies: accessPolicies
    tags: tags
    sku: {
      family: 'A'
      name: 'standard'
    }
  }
}

module cosmosMod 'modules/cosmosdb.bicep' = {
  name: 'cosmosDbDeployment'
  params: {
    name: cosmosDbName
    tags: tags
  }
}

module sbMod 'modules/sbus.bicep' = {
  name: 'serviceBusDeployment'
  params: {
    name: serviceBusName
    skuName: 'Standard'
    tags: tags
  }
}

module vnetMod 'modules/vnet.bicep' = {
  name: 'vnetDeployment'
  params: {
    name: vnetName
    tags: tags
    addressPrefix: '10.0.0.0/16'
    subnets: [
      {
        name: 'aks-sys-subnet'
        addressPrefix: '10.0.0.0/24'
      }
      {
        name: 'aks-usr-subnet'
        addressPrefix: '10.0.1.0/24'
      }
    ]
  }
}

module aksMod 'modules/aks.bicep' = {
  name: aksName
  params: {
    name: aksName
    addOns: {}
    aksVersion: '1.20.7'
    networkPlugin: 'azure'
    enableAutoScaling: true
    aksAgentOsDiskSizeGB: 30
    adminGroupObjectID: adminGroupId
    aksSystemSubnetId: vnetMod.outputs.subnets[0].id
    aksUserSubnetId: vnetMod.outputs.subnets[1].id
    linuxAdminUserName: 'localadmin'
    logAnalyticsWorkspaceId: wksMod.outputs.workspaceId
    sshPublicKey: sshPublicKey
    tags: tags
  }
}

resource acr 'Microsoft.ContainerRegistry/registries@2020-11-01-preview' existing = {
  name: acrName
}

resource AssignAcrPullToAks 'Microsoft.Authorization/roleAssignments@2021-04-01-preview' = {
  name: guid(resourceGroup().id, acrName, 'AssignAcrPullToAks')
  scope: acr
  properties: {
    description: 'Assign AcrPull role to AKS'
    principalId: aksMod.outputs.aksKubeletPrincipalId
    principalType: 'ServicePrincipal'
    roleDefinitionId: acrPullRoleId
  }
}

output cosmosDbId string = cosmosMod.outputs.id
output acrId string = acr.id
output kvId string = kvMod.outputs.id
output kvUri string = kvMod.outputs.keyVaultUri
