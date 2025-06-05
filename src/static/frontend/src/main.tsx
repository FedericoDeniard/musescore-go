import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.tsx";
import "./index.css";
import { Toaster } from "sonner";

import { Amplify } from "aws-amplify";
import config from "./aws-exports.js";
Amplify.configure(config);

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <Toaster position="bottom-center" richColors={true} />
    <App />
  </React.StrictMode>
);
