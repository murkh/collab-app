// src/collab/Editor.tsx
import React, { useEffect, useRef, useState } from "react";
import * as Y from "yjs";
import { HocuspocusProvider } from "@hocuspocus/provider";
import { Controlled as CodeMirror } from "react-codemirror2";

type User = {
  clientId: number;
  name: string;
  color: string;
};

export default function Editor({
  docId,
  token,
}: {
  docId: string;
  token: string;
}) {
  const [value, setValue] = useState("");
  const [users, setUsers] = useState<User[]>([]);
  const ydocRef = useRef<Y.Doc>(null);
  const providerRef = useRef<HocuspocusProvider>(null);

  useEffect(() => {
    const ydoc = new Y.Doc();

    // Setup provider with JWT auth
    const provider = new HocuspocusProvider({
      url: "wss://yourdomain.com/collab",
      name: docId,
      document: ydoc,
      token: token,
    });

    const ytext = ydoc.getText("codemirror");

    // Bind Y.Text to state
    ytext.observe(() => {
      setValue(ytext.toString());
    });

    // --- Awareness setup ---
    const awareness = provider.awareness;

    // Example: assign random username + color (replace with real user info)
    const userName = "User-" + Math.floor(Math.random() * 1000);
    const userColor = "#" + Math.floor(Math.random() * 16777215).toString(16);

    awareness?.setLocalStateField("user", {
      name: userName,
      color: userColor,
    });

    awareness?.on("change", () => {
      const states: User[] = [];
      awareness.getStates().forEach((state: any, clientId: number) => {
        if (state.user) {
          states.push({
            clientId,
            name: state.user.name,
            color: state.user.color,
          });
        }
      });
      setUsers(states);
    });

    ydocRef.current = ydoc;
    providerRef.current = provider;

    return () => {
      provider.destroy();
      ydoc.destroy();
    };
  }, [docId, token]);

  return (
    <div className="flex flex-col h-screen">
      {/* Online users */}
      <div className="flex space-x-2 p-2 border-b">
        {users.map((u) => (
          <div
            key={u.clientId}
            className="px-2 py-1 rounded text-white"
            style={{ backgroundColor: u.color }}
          >
            {u.name}
          </div>
        ))}
      </div>

      {/* Editor */}
      <div className="flex-1">
        <CodeMirror
          value={value}
          options={{
            mode: "javascript",
            theme: "material",
            lineNumbers: true,
          }}
          onBeforeChange={(_editor, _data, val) => {
            const ytext = ydocRef.current?.getText("codemirror");
            if (ytext) {
              ytext.delete(0, ytext.length);
              ytext.insert(0, val);
            }
          }}
        />
      </div>
    </div>
  );
}
