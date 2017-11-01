# Panthalassa
> A Javascript + Flow implementation of panthalassa

## Api
> Panthalassa is under heavy development. Things will change fast.

### Integration
1. Install the project
2. Use `reqiure('Panthalassa')` to include it
3. Install the mobile secure storage

### Methods

#### ETH
* `eth.createPrivateKey() : Promise<string>`
    * Parameter: - 
    * Returns: A promise that resolves with a string

* `eth.savePrivateKey(privateKey: string, pw: ?string, pwConfirm: ?string) : Promise<void>`
    * Parameter:
        * privateKey, should be the private key in hex form.
        * pw, OPTIONAL. The password which will be used for encryption.
        * pw, OPTIONAL. The password which will be used for encryption.
    * Returns: Void promise

* `eth.allKeyPairs() : Promise<Array<{key:string, value:any}>>`
    * Parameter: -
    * Returns: Array of key value objects
        ````js
          [
              {
                  key: '0x2F3D5824C04cc1ABbC070568860A8f7b838b1cab',
                  value: {
                      encryption: 'AES-256',
                      value: 'ac750146531db743fdfb71d83d08ea8cd66b1f9aa24ebc42184f2c33955a9bd5',
                      encrypted: false,
                      version: '1.0.0'
                  }
              }
          ]
        ````    
* `eth.getPrivateKey(address:string) : Promise<{}>`
    * Parameter:
        * address: Is an ethereum address
    * Returns: Promise that resolves with one key value pair. The key value pair will look like this (the key prop is the public address and the value is the private key + information about encryption etc): 
    ````js
      //PrivateKey
      {
          key: '0x2F3D5824C04cc1ABbC070568860A8f7b838b1cab',
          value: {
              encryption: 'AES-256',
              value: 'ac750146531db743fdfb71d83d08ea8cd66b1f9aa24ebc42184f2c33955a9bd5',
              encrypted: false,
              version: '1.0.0'
          }
      }
    ````
* `eth.deletePrivateKey(address:string) : Promise<void}>`
    * Parameter:
        * address: Is an ethereum address
    * Returns: Void promise
    
* `eth.decryptPrivateKey(privateKey: object, reason:string, topic:string) : Promise<void>`
    * Parameter: 
        * privateKey: Is a object that contains a key(ethereum address) and an object as value
          ````js
          //PrivateKey
          {
              key: '0x2F3D5824C04cc1ABbC070568860A8f7b838b1cab',
              value: {
                  encryption: 'AES-256',
                  value: 'ac750146531db743fdfb71d83d08ea8cd66b1f9aa24ebc42184f2c33955a9bd5',
                  encrypted: false,
                  version: '1.0.0'
              }
          }
          ````
        * reason: This string can be something like "Encrypt you private key to display it". It is used for the alert. 
        * topic: Can be something like "Sign transaction"
    * Response: The response will be a Promise that resolves with the raw private key BUT you need to subscribe to the `eth:decrypt-private-key` event in order to be able to resolve the promise. Read more about it in the event section.
### Events

* `eth:decrypt-private-key`
    * Why? You may ask why there is need for this event. The thing is, that Panthalassa does some background work and need a way to tell pangea "Hey, we are using the private key to do some stuff, please show an alert and ask the user to decrypt his private key". Let me explain how this event relates to the `eth.decryptPrivateKey` method. When you call the `decryptPrivateKey` method a Promise will be returned. The ONLY way to resolve the promise is, to call the event "successor" function with the correct password. When you call the successor with the correct password the promise returned by `eth.decryptPrivateKey` will be resolved as well as the promise returned by the "successor" callback. If you have any question's about this feel free to open an issue.
    * Parameter: The parameters will arrive in an object
        * successor: callback
            * Parameter: 
                * password: string
            * Return: Promise that will be rejected or resolved
        * killer: callback
            * Parameter: -
            * Return: -
        * reason: string
        * topic: string
## FAQ

**I heard this is supposed to be the backend of pangea, how can the be?**
>ok, so your backend is not a common backend where you make a few http request, get some data back and done. Our backend is a decentraliced meshnetwork (it will be in some days). That means, each device in the network is a "server" (not really a server but a node). There for this need's to run on each device (like smartphones and laptops). The owner of the device will be able to communicate with other people in the network since he is a node in it.

## Specification

#### Secure Storage
> The secure storage is used to save critical information such as private keys in a save environment. 

You can find the specification [here](./src/specification/secureStorageInterface.js)