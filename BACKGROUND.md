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

These six (6) parties have four (4) different permissioned sets of data. For example, data among the farmer, distributor, and shipper is shared in a common set, and the distributor also shares a set with the wholesaler. However, the wholesaler is not able to see any data within the farmer-distributor-shipper set.  In other tools, these permission sets are called 'channels' so we'll do the same for the remainder of this document.

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

### Considerations/Requirements
1. The number of blockchains used can be flexible such that it is not administratively burdensome. It should have the ability to operate completely within a private setting (no public chains), but also operate on a public chain.  See https://www.coindesk.com/microsoft-ey-and-consensys-present-new-way-for-big-biz-to-use-public-ethereum
2. The possible permutations of channels is a [power set](https://en.wikipedia.org/wiki/Power_set), meaning where *s* is the number of parties, the number of potential data access sets is *2^s*. Our avocado supply chain may have 64 channels.
3. The solution may operate as a "layer 2".  The application will likely be fed by a message bus system that already integrates with Enterprise legacy systems.
4. What happens when channel membership changes (revocation) to either add or remove a party? 
5. The solution should support metadata-visible and metadata-hidden models.  Currently, chappe is a metadata-hidden because each node on the network has the same private key to the message bus contract, so parties cannot deduce who is posting messages.  This could be addressed by changing the keys per party, etc. 
6. Recipients of messages should also post acknowledgements of initial messages, which confirms receipt. This acknowledgement should have cryptographic proof that the recipient was able to receive and decrypt the message.  
7. In Chappe's current design, participants do not know how many other participants are using the network or even how many other participants that are in channels in which they operate. This could be a problem, so perhaps we could require unencrypted messages periodically for "proof of keys" - I imagine this will vary by implementation.

### Extra Credit
1. Implementing a protocol such as JEDI (https://arxiv.org/pdf/1905.13369.pdf) would allow for:
    - Hierarchy of permissions (delegation)
    - Expiration of Permissions
2. How does initial channel key distribution work?