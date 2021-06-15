import { useCallback } from 'react';
import Head from 'next/head'
import Image from 'next/image'
import csvParse from 'csv-parse/lib/sync'
import { useAtom } from 'jotai';
import { setItemStateAtom, clearItemsAtom } from '../store';
import styles from '../styles/Home.module.css'

export default function Items(): React.FC {
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
    </li>
  )
};

function getClassFromState(state) {
  if(state === 'pending') {
    return 'animation-pulse';
  }
  if(state === 'done') {
    return 'line-through';
  }
  return '';
}
