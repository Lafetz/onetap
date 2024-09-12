import { Server } from "socket.io";

import http from "http";

import redisClient from "../redis";
const redisSubscriber = redisClient.duplicate();
const socketServer = async (
  server: http.Server<typeof http.IncomingMessage, typeof http.ServerResponse>
) => {
  const io = new Server(server);
  await redisSubscriber.connect();
  const connectedUsers: Record<string, string> = {};
  redisSubscriber.subscribe("notification", (message) => {
    console.log("got notification");
    const notification = JSON.parse(message);
    const { CustomerID, Desc } = notification;
    const socketId = connectedUsers[CustomerID];
    if (socketId) {
      io.to(socketId).emit("notification", Desc);
    }
  });
  io.on("connection", (socket) => {
    socket.on("id", (CustomerID: string) => {
      console.log("got id");
      if (typeof CustomerID !== "string") {
        console.error("Invalid CustomerID. It must be a string.");
        return;
      }
      connectedUsers[CustomerID] = socket.id;
      socket.on("disconnect", () => {
        delete connectedUsers[CustomerID];
      });
    });
  });
};
export default socketServer;
