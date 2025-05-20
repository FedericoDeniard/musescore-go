import { useState } from "react";
import "./index.css";
import { Loader } from "../Loader";

interface SearchBoxProps {
  onClick: (url: string) => void;
  loading: boolean;
}

export const SearchBox = ({ onClick, loading }: SearchBoxProps) => {
  const [url, setUrl] = useState("");

  const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUrl(e.target.value);
  };

  const handleButtonClick = () => {
    onClick(url);
    setUrl("");
  };

  return (
    <div className="searchBox">
      <input
        className="search"
        type="text"
        placeholder="Enter the musescore URL"
        value={url}
        onChange={handleInput}
        disabled={loading}
      />
      {!loading ? (
        <input
          className="button"
          type="button"
          value="Download"
          onClick={handleButtonClick}
        />
      ) : (
        <Loader />
      )}
    </div>
  );
};
