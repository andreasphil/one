import { newNavigationAction, useCommands } from "../lib/commands.js";
import { getIcons } from "../state/icons.js";
import { getNotes } from "../state/notes.js";
import { getTags } from "../state/tags.js";

export function init() {
  const { getElement: icon } = getIcons();
  const { notes } = getNotes();
  const { tags } = getTags();
  const { register } = useCommands();

  /** @type {import("../lib/commands.js").Command[]} */
  const notesCommands = notes.map((i) => ({
    id: `note:open:${i.Slug}`,
    name: i.Title,
    groupName: "Notes",
    icon: icon("StickyNote"),
    action: newNavigationAction(`/notes/${i.Slug}/`),
  }));

  /** @type {import("../lib/commands.js").Command[]} */
  const tagsCommands = tags.map((i) => ({
    id: `tag:open:${i}`,
    name: i,
    alias: [`#${i}`],
    groupName: "Tags",
    icon: icon("Tag"),
    action: newNavigationAction(`/tags/${i}/`),
  }));

  /** @type {import("../lib/commands.js").Command[]} */
  const commands = [
    {
      id: "open:search",
      name: "Search",
      chord: "gs",
      groupName: "Open",
      icon: icon("Search"),
      action: newNavigationAction("/search/"),
    },
    {
      id: "open:notes",
      name: "Notes",
      chord: "gn",
      groupName: "Open",
      icon: icon("StickyNote"),
      action: newNavigationAction("/notes/")
    },
    {
      id: "open:today",
      name: "Today",
      chord: "gt",
      groupName: "Open",
      icon: icon("Calendar"),
      action: newNavigationAction(`/notes/${new Date().toISOString().substring(0, 10)}/`),
    },
  ];

  register(...tagsCommands, ...notesCommands, ...commands);
}
