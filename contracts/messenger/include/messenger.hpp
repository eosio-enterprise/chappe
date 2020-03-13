#include <eosio/eosio.hpp>
using namespace eosio;
using std::string;

CONTRACT messenger : public contract {
   public:
      using contract::contract;
    
      ACTION pub( string ipfs_hash, string memo );

      using pub_action = action_wrapper<"pub"_n, &messenger::pub>;
};