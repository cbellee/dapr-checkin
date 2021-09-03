param name string
param tags object

resource cs 'Microsoft.DocumentDB/databaseAccounts@2021-06-15' = {
  location: resourceGroup().location
  kind: 'GlobalDocumentDB'
  name: name
  tags: tags
  properties: {
    locations: [
      {
        locationName: resourceGroup().location
        isZoneRedundant: false
      }
    ]
    databaseAccountOfferType: 'Standard'
    createMode: 'Default'
    publicNetworkAccess: 'Enabled'
  }
}

output name string = cs.name
output id string = cs.id
