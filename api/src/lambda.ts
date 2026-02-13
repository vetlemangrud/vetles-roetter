import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { handle } from "hono/aws-lambda";
import { carrotRoutes } from "./routes/carrots.js";

const app = new Hono();

app.use("*", logger());
app.use("*", cors());

app.get("/", (c) => c.json({ status: "ok", service: "gulrot-tracker" }));

app.route("/api/carrots", carrotRoutes);

export const handler = handle(app);
