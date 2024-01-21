module test1::my_module {

    use sui::object::{Self, UID};
    use sui::transfer;
    use sui::tx_context::{Self, TxContext};

    struct Sword has key, store {
        id: UID,
        magic: u64,
        strength: u64,
    }

    struct Forge has key, store {
        id: UID,
        swords_created: u64,
    }

    // Module initializer to be executed when this module is published
    fun init(ctx: &mut TxContext) {
        let admin = Forge {
            id: object::new(ctx),
            swords_created: 0,
        };
        // Transfer the forge object to the module/package publisher
        transfer::transfer(admin, tx_context::sender(ctx));
    }

    // Accessors required to read the struct attributes
    public fun magic(self: &Sword): u64 {
        self.magic
    }

    public fun strength(self: &Sword): u64 {
        self.strength
    }

    public fun swords_created(self: &Forge): u64 {
        self.swords_created
    }

    public fun mint_to_sender(self: &mut Forge, magic: u64, strength: u64, ctx: &mut TxContext) {
        // create a sword
        let sword = Sword {
            id: object::new(ctx),
            magic: magic,
            strength: strength,
        };

        self.swords_created = self.swords_created + 1;

        // transfer the sword to sender
        transfer::public_transfer(sword, tx_context::sender(ctx));
    }

    public fun sword_transfer(sword: Sword, recipient: address, _ctx: &mut TxContext) {
        // transfer the sword
        transfer::public_transfer(sword, recipient);
    } 
}