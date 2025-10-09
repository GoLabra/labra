type Document = Record<string, any>;

function descendingComparator(a: Document, b: Document, sortBy: string): number {
  // When compared to something undefined, always returns false.
  // This means that if a field does not exist from either element ('a' or 'b') the return will be 0.

  if (b[sortBy]! < a[sortBy]!) {
    return -1;
  }

  if (b[sortBy]! > a[sortBy]!) {
    return 1;
  }

  return 0;
}

function getComparator(sortDir: string, sortBy: string) {
  return (
    sortDir === 'desc'
      ? (a: Document, b: Document) => descendingComparator(a, b, sortBy)
      : (a: Document, b: Document) => -descendingComparator(a, b, sortBy)
  );
}

export function applySort<T = Document>(
  documents: T[],
  sortBy: string,
  sortDir: 'asc' | 'desc'
): T[] {
  const comparator = getComparator(sortDir, sortBy);
  const stabilizedThis = documents.map((el, index) => [el, index]);

  stabilizedThis.sort((a, b) => {
    // @ts-ignore
    const newOrder = comparator(a[0], b[0]);

    if (newOrder !== 0) {
      return newOrder;
    }

    // @ts-ignore
    return a[1] - b[1];
  });

  // @ts-ignore
  return stabilizedThis.map((el) => el[0]);
}
