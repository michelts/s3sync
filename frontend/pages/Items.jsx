import { useCallback } from 'react';
import Head from 'next/head'
import Image from 'next/image'
import csvParse from 'csv-parse/lib/sync'
import { useAtom } from 'jotai';
import { setItemStateAtom, clearItemsAtom } from '../store';
import styles from '../styles/Home.module.css'

export default function Items() {
  const [items, setItemState] = useAtom(setItemStateAtom);
  const [, clearItems] = useAtom(clearItemsAtom);

  return (
    <div>
      <h1 className="mb-3">
        <span className="text-xl">
          Synchonization ongoing
        </span>
        {' '}
        <button
          className="text-blue-600 hover:underline"
          onClick={clearItems}
        >(cancel)</button>
      </h1>
      <ul>
        {items.map((item) => <Item item={item} />)}
      </ul>
    </div>
  )
}

const Item = ({ item }) => {
  const [publisher, publication, issue, state] = item;
  return (
    <li
      key={issue}
      className={getClassFromState(state)}
    >
      {`publisher/${publisher}/publication/${publication}/issue/${issue}/`}
      {state === 'pending' && <Spin />}
    </li>
  )
};

function getClassFromState(state) {
  if(state === 'done') {
    return 'line-through';
  }
  return '';
}

const Spin = () => (
  <svg className="animate-spin ml-1 mr-3 h-4 w-4 text-indigo-500 inline-block" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
    <circle className="opacity-50" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
    <path className="opacity-100" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
  </svg>
)
