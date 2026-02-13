import { sqliteTable, text, integer } from "drizzle-orm/sqlite-core";

export const carrots = sqliteTable("carrots", {
  id: integer("id").primaryKey({ autoIncrement: true }),
  eatenAt: integer("eaten_at", { mode: "timestamp" })
    .notNull()
    .$defaultFn(() => new Date()),
});

export type Carrot = typeof carrots.$inferSelect;
export type NewCarrot = typeof carrots.$inferInsert;
