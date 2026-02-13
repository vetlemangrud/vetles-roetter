import { Hono } from "hono";
import { db } from "../db/index.js";
import { carrots } from "../db/schema.js";
import { eq, desc, sql, gte } from "drizzle-orm";

export const carrotRoutes = new Hono()
  .post("/", async (c) => {
    const result = await db.insert(carrots).values({}).returning();
    return c.json(result[0], 201);
  })

  .get("/", async (c) => {
    const limit = Number(c.req.query("limit")) || 50;
    const offset = Number(c.req.query("offset")) || 0;

    const result = await db
      .select()
      .from(carrots)
      .orderBy(desc(carrots.eatenAt))
      .limit(limit)
      .offset(offset);

    return c.json(result);
  })

  .get("/stats", async (c) => {
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

    return c.json({
      total: total.count,
      today: today.count,
      thisWeek: thisWeek.count,
    });
  })

  .delete("/:id", async (c) => {
    const id = Number(c.req.param("id"));

    const result = await db
      .delete(carrots)
      .where(eq(carrots.id, id))
      .returning();

    if (result.length === 0) {
      return c.json({ error: "Carrot not found" }, 404);
    }

    return c.json({ deleted: result[0] });
  });
