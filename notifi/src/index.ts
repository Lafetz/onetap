import express from "express";
import http from "http";
import socketServer from "./socket/socket.server";
//

import { createClient } from "redis";

// Retrieve Redis credentials from environment variables
const REDIS_URL = process.env.REDIS_URL || "cache:6379";
const REDIS_PASS = process.env.REDIS_PASS || "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81";

// Create a Redis client
const redisClient = createClient({
  url: `redis://${REDIS_URL}`,
  password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
});

const redisSubscriber = redisClient.duplicate(); // Duplicate for subscribing

const connectToRedis = async () => {
  try {
    await redisClient.connect();
    await redisSubscriber.connect();
    console.log("Connected to Redis");
  } catch (err) {
    console.error("Redis connection error:", err);
  }
};

// Call the function to connect to Redis
connectToRedis();
//
const app = express();
const server = http.createServer(app);
socketServer(server);
server.listen(4000, () => {
  console.log("server started");
});
