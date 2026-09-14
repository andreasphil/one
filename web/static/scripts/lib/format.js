const dateFormatter = new Intl.DateTimeFormat("de", { dateStyle: "medium", timeStyle: undefined });

/** @param {Date} input */
export function formatDate(input) {
  return dateFormatter.format(input);
}
