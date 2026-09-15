/** Client-side types for note metadata. Should mirror noteMeta in mapper.go. */
export type NoteMeta = {
  Title: string;
  Slug: string;
  Date?: string;
  IsDailyNote: boolean;
  IsChildNote: boolean;
};
