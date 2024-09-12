import { createClient } from "redis";

// Retrieve Redis credentials from environment variables
const REDIS_URL = process.env.REDIS_URL || "localhost:6379";
const REDIS_PASS = process.env.REDIS_PASS || "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81";
const redisConnection = createClient({
  url: `redis://${REDIS_URL}`,
  password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
});

const redisSubscriber = redisConnection.duplicate(); // Duplicate for subscribing

const connectToRedis = async () => {
  try {
    await redisConnection.connect();
    await redisSubscriber.connect();
    console.log("Connected to Redis");
  } catch (err) {
    console.error("Redis connection error:", err);
  }
};

// Call the function to connect to Redis
connectToRedis();
export default redisConnection;
