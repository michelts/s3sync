import { useCallback } from 'react';
import Head from 'next/head'
import Image from 'next/image'
import csvParse from 'csv-parse/lib/sync'
import { useAtom } from 'jotai'
import { itemsAtom } from '../store'
import CSVImport from './CSVImport'
import Items from './Items'

export default function Home() {
  const [items] = useAtom(itemsAtom);

  if(!items?.length) {
    return <CSVImport />
  }

  return <Items />
}
