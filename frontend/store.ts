import { atom } from 'jotai'

export const itemsAtom = atom([]);

export const setItemsAtom = atom(
  null,
  async (get, set, values) => {
    set(itemsAtom, values)
    for(let index=0; index<values.length; index++) {
      const [publisher, publication, issue] = values[index];
      set(
        itemsAtom,
        values.map((item, itemIndex) => (itemIndex === index) ? [...item.slice(0, 2), 'pending'] : item),
      )
      await new Promise(resolve => setTimeout(resolve, 300))
      set(
        itemsAtom,
        values.map((item, itemIndex) => (itemIndex === index) ? [...item.slice(0, 2), 'done'] : item),
      )
    }
  }
);

export const clearItemsAtom = atom(
  null,
  (get, set) => set(itemsAtom, []),
);

export const setItemStateAtom = atom(
  (get) => get(itemsAtom),
  (get, set, { issue, state }) => set(
    itemsAtom,
    get(itemsAtom).map((item) => (
      (item[2] === issue)
        ? [...item.slice(0, 3), state]
        : item
    ))
  )
);
