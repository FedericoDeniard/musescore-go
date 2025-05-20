import { useEffect, useState } from "react";

import { KEYS } from "./constants/keys";

import "./App.css";
import { Title } from "./components/Title";
import { SearchBox } from "./components/SearchBox";
import { toast } from "sonner";

function App() {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    const downloadUrl = async () => {
      if (url) {
        setLoading(true);

        toast.promise(
          fetch(KEYS.url, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
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
    <div className="container-title">
      <Title />
      <SearchBox
        onClick={(getUrl: string) => setUrl(getUrl)}
        loading={loading}
      />
    </div>
  );
}

export default App;
