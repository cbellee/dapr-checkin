param name string
param skuName string = 'Standard'
param tags object

resource sb 'Microsoft.ServiceBus/namespaces@2021-06-01-preview' = {
  name: name
  tags: tags
  location: resourceGroup().location
  sku: {
    name: skuName
  }
}

output id string = sb.id
output endpoint string = sb.properties.serviceBusEndpoint
