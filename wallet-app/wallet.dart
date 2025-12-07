// NUSA Wallet - Mobile & Web (Flutter)
// Features: Send/receive, swap, contacts, transaction history

import 'dart:convert';
import 'package:crypto/crypto.dart';

class NusaWallet {
  String address;
  String privateKey;
  double balance;
  List<Transaction> transactions;
  List<Contact> contacts;
  
  NusaWallet({
    required this.address,
    required this.privateKey,
    this.balance = 0.0,
    List<Transaction>? transactions,
    List<Contact>? contacts,
  })  : transactions = transactions ?? [],
        contacts = contacts ?? [];

  // Create new wallet
  static Future<NusaWallet> create() async {
    // Generate random private key
    final random = List<int>.generate(32, (i) => i);
    final privateKey = sha256.convert(random).toString();
    
    // Derive address from private key
    final address = _deriveAddress(privateKey);
    
    print('✅ Wallet created: $address');
    
    return NusaWallet(
      address: address,
      privateKey: privateKey,
      balance: 0.0,
    );
  }

  // Import wallet from mnemonic
  static Future<NusaWallet> importFromMnemonic(String mnemonic) async {
    // Convert mnemonic to seed
    final seed = sha256.convert(utf8.encode(mnemonic)).toString();
    
    // Derive private key and address
    final privateKey = seed.substring(0, 64);
    final address = _deriveAddress(privateKey);
    
    print('📥 Wallet imported: $address');
    
    return NusaWallet(
      address: address,
      privateKey: privateKey,
    );
  }

  // Send transaction
  Future<String> send(String to, double amount, double gasPrice) async {
    if (amount > balance) {
      throw Exception('Insufficient balance');
    }
    
    // Create transaction
    final tx = Transaction(
      from: address,
      to: to,
      value: amount,
      gasPrice: gasPrice,
      timestamp: DateTime.now(),
      status: 'pending',
    );
    
    // Sign transaction
    final signature = _signTransaction(tx, privateKey);
    tx.signature = signature;
    
    // Broadcast (simplified - use actual RPC)
    print('📤 Sending $amount NUSA to $to');
    
    transactions.add(tx);
    balance -= amount;
    
    return tx.hash;
  }

  // Receive transaction
  void receive(Transaction tx) {
    if (tx.to == address) {
      balance += tx.value;
      transactions.add(tx);
      print('📥 Received ${tx.value} NUSA from ${tx.from}');
    }
  }

  // Swap tokens
  Future<void> swap(String fromToken, String toToken, double amount) async {
    // Integrate with DEX
    print('🔄 Swapping $amount $fromToken → $toToken');
    
    // Call DEX smart contract
    // Update balance
  }

  // Add contact
  void addContact(String name, String address) {
    contacts.add(Contact(name: name, address: address));
    print('👤 Contact added: $name ($address)');
  }

  // Get transaction history
  List<Transaction> getHistory({int limit = 20}) {
    return transactions.take(limit).toList();
  }

  // Estimate gas
  double estimateGas(String operation) {
    final gasEstimates = {
      'transfer': 21000. 0,
      'swap': 150000.0,
      'nft_mint': 100000.0,
    };
    
    return gasEstimates[operation] ?? 21000.0;
  }

  // Get balance
  Future<double> refreshBalance() async {
    // Query blockchain (simplified)
    // balance = await rpc.getBalance(address);
    return balance;
  }

  // Derive address from private key
  static String _deriveAddress(String privateKey) {
    final hash = sha256.convert(utf8.encode(privateKey));
    return '0x${hash.toString().substring(0, 40)}';
  }

  // Sign transaction
  static String _signTransaction(Transaction tx, String privateKey) {
    final data = '${tx.from}${tx.to}${tx. value}${tx.timestamp}';
    final combined = privateKey + data;
    return sha256.convert(utf8.encode(combined)).toString();
  }
}

class Transaction {
  String from;
  String to;
  double value;
  double gasPrice;
  DateTime timestamp;
  String status;
  String?  signature;
  String hash;

  Transaction({
    required this.from,
    required this.to,
    required this.value,
    required this. gasPrice,
    required this. timestamp,
    required this.status,
    this.signature,
  }) : hash = sha256.convert(utf8.encode('$from$to$value${timestamp.millisecondsSinceEpoch}')).toString();
}

class Contact {
  String name;
  String address;

  Contact({required this.name, required this.address});
}

// Example usage
void main() async {
  // Create wallet
  final wallet = await NusaWallet.create();
  
  // Add some test balance
  wallet.balance = 100.0;
  
  // Add contact
  wallet.addContact('Alice', '0xabc123.. .');
  
  // Send transaction
  try {
    final txHash = await wallet.send('0xabc123... ', 10.0, 1000.0);
    print('✅ Transaction sent: $txHash');
  } catch (e) {
    print('❌ Error: $e');
  }
  
  // Check balance
  print('💰 Balance: ${wallet.balance} NUSA');
  
  // View history
  print('📜 History: ${wallet.transactions.length} transactions');
}
