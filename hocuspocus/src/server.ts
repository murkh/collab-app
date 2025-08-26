import { Hocuspocus, Server } from "@hocuspocus/server";
import { Logger } from "@hocuspocus/extension-logger";
import { Redis } from "@hocuspocus/extension-redis";
import { Database } from "@hocuspocus/extension-database";
import jwt from "jsonwebtoken";
import { Pool } from "pg";

const REDIS_HOST = process.env.REDIS_HOST || "localhost";
const REDIS_PORT = process.env.REDIS_PORT || "6379";
const PG_DSN =
  process.env.PG_DSN || "postgres://postgres:postgres@localhost:5432/postgres";
const JWT_PUB = process.env.JWT_PUB || "";

const pg = new Pool({ connectionString: PG_DSN });

const server = new Server({
  port: process.env.PORT ? Number(process.env.PORT) : 1234,
  extensions: [
    new Logger(),
    new Redis({
      host: "127.0.0.1",
      port: Number(REDIS_PORT),
    }),
    new Database({
      fetch: ({ documentName }) => {
        return new Promise((resolve, reject) => {
          pg.query(
            `SELECT state, version from yjs_snapshots WHERE doc_id = $1
                  ORDER BY version DESC LIMIT 1`,
            [documentName]
          )
            .then((res) => {
              if (res.rowCount === 0) return resolve(null);
              resolve(res.rows[0].state);
            })
            .catch(reject);
        });
      },
      store: async ({ documentName, state }) => {
        const ver = Date.now();
        await pg.query(
          `INSERT INTO yjs_snapshots(doc_id, version, state, created_at) values($1, $2, $3, now())`,
          [documentName, ver, state]
        );
      },
    }),
  ],
  async onAuthenticate({ token }) {
    if (!token) return false;
    try {
      const payload = jwt.verify(token as string, JWT_PUB, {
        algorithms: ["RS256"],
      });
      return payload;
    } catch (e) {
      console.error("auth failed", e);
      return false;
    }
  },
});

server.listen();
console.log("Hocuspocus running");
