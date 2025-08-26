import { useEffect, useState } from "react";
import Editor from "./collab/Editor";

function App() {
  const [token, setToken] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/collab/token?docId=doc-123", {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => setToken(data.token));
  }, []);

  if (!token) return <div>Loading…</div>;

  return <Editor docId="doc-123" token={token} />;
}

export default App;
