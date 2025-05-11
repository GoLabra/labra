type Document = any;

export function applyPagination<T = Document>(
  documents: T[],
  page: number,
  rowsPerPage: number
): T[] {
  return documents.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage);
}
