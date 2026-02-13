#!/usr/bin/env node

import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { drizzle } from "drizzle-orm/libsql";
import { createClient } from "@libsql/client";
import { sqliteTable, text, integer } from "drizzle-orm/sqlite-core";
import { eq, desc, sql, gte } from "drizzle-orm";
import { z } from "zod";

const carrots = sqliteTable("carrots", {
  id: integer("id").primaryKey({ autoIncrement: true }),
  eatenAt: integer("eaten_at", { mode: "timestamp" })
    .notNull()
    .$defaultFn(() => new Date()),
});

const client = createClient({
  url: process.env.TURSO_DATABASE_URL || "file:local.db",
  authToken: process.env.TURSO_AUTH_TOKEN,
});

const db = drizzle(client, { schema: { carrots } });

const server = new McpServer({
  name: "gulrot-tracker",
  version: "0.1.0",
});

server.tool(
  "log_carrot",
  "Log that you ate a carrot. Records current timestamp.",
  {},
  async () => {
    const result = await db.insert(carrots).values({}).returning();
    const carrot = result[0];
    return {
      content: [
        {
          type: "text",
          text: `Logged carrot #${carrot.id} at ${carrot.eatenAt.toISOString()}`,
        },
      ],
    };
  }
);

server.tool(
  "list_carrots",
  "List recent carrot entries.",
  {
    limit: z.number().optional().default(10).describe("Number of entries to return (default 10)"),
  },
  async ({ limit }) => {
    const result = await db
      .select()
      .from(carrots)
      .orderBy(desc(carrots.eatenAt))
      .limit(limit);

    const entries = result.map(
      (c) => `#${c.id}: ${c.eatenAt.toISOString()}`
    );

    return {
      content: [
        {
          type: "text",
          text: entries.length > 0 ? entries.join("\n") : "No carrots logged yet.",
        },
      ],
    };
  }
);

server.tool(
  "get_stats",
  "Get carrot consumption statistics: today, this week, and total.",
  {},
  async () => {
    const now = new Date();
    const todayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const weekStart = new Date(todayStart);
    weekStart.setDate(weekStart.getDate() - 7);

    const [total] = await db
      .select({ count: sql<number>`count(*)` })
      .from(carrots);

    const [today] = await db
      .select({ count: sql<number>`count(*)` })
      .from(carrots)
      .where(gte(carrots.eatenAt, todayStart));

    const [thisWeek] = await db
      .select({ count: sql<number>`count(*)` })
      .from(carrots)
      .where(gte(carrots.eatenAt, weekStart));

    return {
      content: [
        {
          type: "text",
          text: `Carrot Stats:\n- Today: ${today.count}\n- This week: ${thisWeek.count}\n- Total: ${total.count}`,
        },
      ],
    };
  }
);

server.tool(
  "delete_carrot",
  "Delete a carrot entry by ID.",
  {
    id: z.number().describe("The carrot entry ID to delete"),
  },
  async ({ id }) => {
    const result = await db
      .delete(carrots)
      .where(eq(carrots.id, id))
      .returning();

    if (result.length === 0) {
      return {
        content: [{ type: "text", text: `Carrot #${id} not found.` }],
        isError: true,
      };
    }

    return {
      content: [
        {
          type: "text",
          text: `Deleted carrot #${id} (was logged at ${result[0].eatenAt.toISOString()})`,
        },
      ],
    };
  }
);

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
}

main().catch(console.error);
