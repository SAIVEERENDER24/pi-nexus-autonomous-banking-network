import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.common.serialization.StringSerializer;
import com.fasterxml.jackson.databind.ObjectMapper; // Add Jackson dependency

import java.util.Properties;

public class TransactionProcessor {
    private KafkaProducer<String, String> producer;
    private static final ObjectMapper objectMapper = new ObjectMapper();

    public TransactionProcessor(String bootstrapServers) {
        Properties props = new Properties();
        props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class.getName());
        
        // Reliability settings
        props.put(ProducerConfig.ACKS_CONFIG, "all");
        
        this.producer = new KafkaProducer<>(props);
    }

    public void processTransaction(Transaction transaction) {
        try {
            // Convert object to JSON String
            String transactionJson = objectMapper.writeValueAsString(transaction);
            
            ProducerRecord<String, String> record = new ProducerRecord<>("transactions", transaction.getCardNumber(), transactionJson);
            
            // Send and forget (or use a callback for error handling)
            producer.send(record, (metadata, exception) -> {
                if (exception != null) {
                    exception.printStackTrace();
                } else {
                    System.out.println("Sent to partition: " + metadata.partition() + " offset: " + metadata.offset());
                }
            });
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    public void close() {
        producer.close();
    }

    public static void main(String[] args) {
        TransactionProcessor processor = new TransactionProcessor("localhost:9092");
        
        Transaction transaction = new Transaction("1234567890", 1000, "USA", "VISA");
        processor.processTransaction(transaction);
        
        // Crucial: Ensure data is flushed and connection closed before exiting
        processor.close();
    }
}

class Transaction {
    private String cardNumber;
    private int amount;
    private String country;
    private String cardType;

    public Transaction(String cardNumber, int amount, String country, String cardType) {
        this.cardNumber = cardNumber;
        this.amount = amount;
        this.country = country;
        this.cardType = cardType;
    }

    // Getters are required for Jackson to "see" the data
    public String getCardNumber() { return cardNumber; }
    public int getAmount() { return amount; }
    public String getCountry() { return country; }
    public String getCardType() { return cardType; }
}