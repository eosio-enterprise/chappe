# eosio Permissioned Data

## Use Case

Multiple parties want to process transactions using a blockchain due to the benefits of security, immutability, and transparency. However, the parties want to be able to share some information with some parties without having that information exposed to all parties.

### Example: Avocado Supply Chain

This avocado supply chain example use case has six (6) parties:
1. Farmer
2. Distributor
3. Shipper
4. Wholesaler
5. Retailer
6. Consumer

These six (6) parties have four (4) different permissioned sets of data. For example, data among the farmer, distributor, and shipper is shared in a common set, and the distributor also shares a set with the wholesaler. However, the wholesaler is not able to see any data within the farmer-distributor-shipper set.  

![](https://docs.telos.kitchen/uploads/upload_e2335ab482ec4156d49b3230ce7b6bc7.png)


#### Example Document
``` json
{
    "purchase_order_number": "3749230AZ-90",
    "vendor_id": "vendor-abc-876",
    "vendor_name": "Farmer ABC",
    "customer_id": "Distributor ABC",
    "ship_to_address": "567 Logistics Ln, Any City, ST, 00000",
    "shipper_id": "Shipper ABC",
    "payment_terms": "NET 30",
    "items": [
        {
            "line_number": 1,
            "item_number": "2390-pwdop-8790",
            "quantity": 34,
            "cost": "3.56 USD"
        },
        {
            "line_number": 2,
            "item_number": "2374-fahop-3290",
            "quantity": 12,
            "cost": "0.74 USD"
        }
    ]
}
```

### Considerations
1. The number of blockchains used can be flexible such that it is not administratively burdensome. It should have the ability to operate completely within a private setting (no public chains), although it is likely beneficial to integrate with public chains in some scenarios.
2. The number of permutations of data access sets is a [power set](https://en.wikipedia.org/wiki/Power_set), meaning where *s* is the number of parties, the number of potential data access sets is *2^s*. Our avocado supply chain may have 64 data access sets.
3. The preference is to rely primarily on enforcing access control at the ***blockchain layer*** and secondarily at the network layer OR the application (off-chain) layer.
4. What happens when an existing data access set is changed to either add or remove a party? 

## Option 1: IPFS private network

IPFS natively supports secure peering by generating a shared secret in the configuration of the node. 

A new IPFS private network would be created for each required permutation of data access sets. In the example above, there would four (4) unique IPFS networks and ten (10) nodes. Parties like the Shipper, who are part of multiple groups will have to run nodes on multiple networks.

### IPFS Networks A-D
We could label the various networks to see how they would map to the entities on the network.
![](https://docs.telos.kitchen/uploads/upload_75918bd50471bfdaf93886ab8fa611c6.png)


Kubernetes IPFS may allow for easy managing of many nodes
https://github.com/ipfs/kubernetes-ipfs

The smart contracts on the network would maintain references to the data objects on IPFS, including the network ID. 

The example network may look something like this. I haven't drawn in all the lines as each network would connect to all its peers. 

![](https://docs.telos.kitchen/uploads/upload_2e5a3b942dcac307fbf04c53e4783aca.png)

The application clients would be able to access the blockchain directly but would need to go through the application server in order to access the permissioned data within IPFS.

![](https://docs.telos.kitchen/uploads/upload_a022304272d715c5e188ec3fb3fbd98f.png)

There would likely be a transformation to an additional database both for caching for the client and use of the information in the rest of the organization.

### Adding or Removing Parties
If a party were added to an existing set, they would just need the secret that was generated at the genesis of that network.

If a party were removed from a set, a new network would need to be created for the remaining parties. Or is it simpler to just share a new key without creating a new network?

### Handling Synchronization
Since nodeos and ipfs will be operating wholly in separate spaces, it is possible that transactions that reference an IPFS hash arrive before the data is added to IPFS.  This should be OK because the transaction can't read it anyway. The application layer would need to identify that scenario and then report back to the network that the data was not received, which may be a problem with the IPFS network or the sender never created the data on IPFS.

## Option 2: Private eosio chain for each data access set
Using this option, we would be running ten (10) eosio chains. 

Instead of adding data to IPFS, the data would be transacted to a private eosio chain. 

#### Benefits/Drawbacks?

## Option 3: Other Distributed Database
The private networks (private data) can run on any distributed database. It seems like IPFS is a pretty good option because it has all of the keys and peering built in. 

This database is actually built on top of IPFS and looks like an option: https://github.com/orbitdb

---

## Appendix A: Writing/retrieving data from IPFS
``` JavaScript
// this is a public IPFS node, but obviously we 
// would be writing private data to the private IPFS node
const ipfs = require("nano-ipfs-store").at("https://ipfs.infura.io:5001");

(async () => {

  const doc = JSON.stringify({
    "purchase_order_number": "3749230AZ-90",
    "vendor_id": "vendor-abc-876",
    "vendor_name": "Farmer ABC",
    "customer_id": "Distributor ABC",
    "ship_to_address": "567 Logistics Ln, Any City, ST, 00000",
    "shipper_id": "Shipper ABC",
    "payment_terms": "NET 30",
    "items": [
      {
        "line_number": 1,
        "item_number": "2390-pwdop-8790",
        "quantity": 34,
        "cost": "3.56 USD"
      },
      {
        "line_number": 2,
        "item_number": "2374-fahop-3290",
        "quantity": 12,
        "cost": "0.74 USD"
      }
    ]
  });

  // add this document to IPFS
  const cid = await ipfs.add(doc);

  // this document/object is saved in IPFS under the following CID
  console.log("IPFS cid:", cid);

  // now we can retrieve this object
  console.log(await ipfs.cat(cid));

})();
```
Here's an example of retrieving the data:
https://ipfs.infura.io/ipfs/bafkreif47eemb2rjap4opfy75btj7wpewvszz3rvhmcaeqc6wygwce5tk4
