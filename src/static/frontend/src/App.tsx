import { useEffect, useState } from "react";

import { KEYS } from "./constants/keys";

import "./App.css";
import { Title } from "./components/Title";
import { SearchBox } from "./components/SearchBox";
import { toast } from "sonner";

import "@aws-amplify/ui-react/styles.css";
import { withAuthenticator, Button, View } from "@aws-amplify/ui-react";

import type { AuthUser } from "aws-amplify/auth";
import { fetchAuthSession } from "aws-amplify/auth";

type AppProps = {
  signOut?: () => void;
  user?: AuthUser;
};

function App({ signOut }: AppProps) {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);

  async function getAuthToken() {
    try {
      const session = await fetchAuthSession();
      const token = session.tokens?.idToken?.toString();
      return token;
    } catch (error) {
      console.error("Error obteniendo token", error);
      return null;
    }
  }

  useEffect(() => {
    const downloadUrl = async () => {
      if (url) {
        setLoading(true);

        const token = await getAuthToken();
        const envUrl = KEYS.url;
        const fetchUrl =
          envUrl === "" ? window.location.origin + "/api/scrap" : envUrl;

        toast.promise(
          fetch(fetchUrl, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              Authorization: `Bearer ${token}`,
            },
            body: JSON.stringify({ url }),
          }).then(async (response) => {
            if (!response.ok) {
              throw new Error("Error descargando la partitura.");
            }

            // Obtener el nombre del archivo desde los encabezados
            const contentDisposition = response.headers.get(
              "Content-Disposition"
            );
            let fileName = "sheet.pdf"; // Valor por defecto si no se encuentra el encabezado

            if (contentDisposition) {
              const match = contentDisposition.match(/filename="(.+)"/);
              if (match && match[1]) {
                fileName = match[1];
              }
            }

            const blob = await response.blob();
            const downloadUrl = window.URL.createObjectURL(blob);

            const a = document.createElement("a");
            a.href = downloadUrl;
            a.download = fileName; // Usar el nombre extraído
            a.click();

            // Liberar el objeto URL
            window.URL.revokeObjectURL(downloadUrl);
          }),
          {
            loading: "Downloading...",
            success: "Downloaded successfully",
            error: (error: Error) =>
              error.message.replace(
                "Error invoking remote method 'download':",
                ""
              ) || `Error downloading the sheet`,
            finally: async () => {
              setLoading(false);
              setUrl("");
            },
          }
        );
      }
    };

    downloadUrl();
  }, [url]);

  return (
    <View>
      <div className="container-title">
        <Title />
        <SearchBox
          onClick={(getUrl: string) => setUrl(getUrl)}
          loading={loading}
        />
        <Button variation="primary" colorTheme="error" onClick={signOut}>
          Cerrar sesión
        </Button>
      </div>
    </View>
  );
}

export default withAuthenticator(App);
