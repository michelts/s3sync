import { useCallback } from 'react';
import Head from 'next/head'
import Image from 'next/image'
import csvParse from 'csv-parse/lib/sync'
import { useAtom } from 'jotai';
import { setItemsAtom } from '../store';
import styles from '../styles/Home.module.css'

export default function CSVImport() {
  const [, setItems] = useAtom(setItemsAtom);
  let reader;

  const handleFileRead = (event) => {
    const parsed = csvParse(reader.result);
    setItems(parsed);
  }

  const handleChange = useCallback((event) => {
    const file = event.target.files[0];
    reader = new FileReader();
    reader.onloadend = handleFileRead;
    reader.readAsText(file);
  }, []);

  return (
    <div>
      <h1
        className="text-xl mb-3"
      >
        Sync Issues
      </h1>
      <div className="mb-2">Please pick a CSV file:</div>
      <input
        type="file"
        name="csv-file"
        onChange={handleChange}
      />
    </div>
  )
}

