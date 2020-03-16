#include <eosio/eosio.hpp>
using namespace eosio;
using std::string;

CONTRACT messenger : public contract {
   public:
      using contract::contract;
    
      ACTION pub( string ipfs_hash, string memo );
      ACTION pubmap ( std::map<string, string> payload);
      // ACTION pubbytesmap ( std::map<string, byte[]> payload);
      
};