import { CommandBar } from "../common/commandBar.js";

document.addEventListener("keyup", (event) => {
  if (event.key === "/") {
    const element = document.querySelector("input[type=search]");

    if (element && element instanceof HTMLInputElement && document.activeElement !== element) {
      event.preventDefault();

      element.focus();
      element.setSelectionRange(0, element.value.length);
    }
  }
});

const iconTemplates = document.querySelectorAll("template[id^='clientState.Icons.']");

const icons = iconTemplates.entries().reduce((all, [, i]) => {
  const name = i.id.substring(i.id.lastIndexOf(".") + 1);
  return all.set(name, i.content.firstElementChild);
}, new Map());

function navigationAction(url) {
  return () => {
    navigation.navigate(url);
  };
}

const notes = JSON.parse(document.getElementById("clientState.NotesMeta").textContent).map((i) => ({
  id: `note:open:${i.Slug}`,
  name: i.Title,
  groupName: "Notes",
  icon: icons.get("StickyNote").cloneNode(true),
  action: navigationAction(`/notes/${i.Slug}/`),
}));

const tags = JSON.parse(document.getElementById("clientState.Tags").textContent).map((i) => ({
  id: `tag:open:${i}`,
  name: i,
  alias: [`#${i}`],
  groupName: "Tags",
  icon: icons.get("Tag").cloneNode(true),
  action: navigationAction(`/tags/${i}/`),
}));

CommandBar.define();

CommandBar.instance.registerCommand(
  ...notes,
  ...tags,
  {
    id: "open:search",
    name: "Search",
    chord: "gs",
    groupName: "Open",
    icon: icons.get("Search").cloneNode(true),
    action: navigationAction("/search/"),
  },
  {
    id: "open:today",
    name: "Today",
    chord: "gt",
    groupName: "Open",
    icon: icons.get("Calendar"),
    action: navigationAction(`/notes/${new Date().toISOString().substring(0, 10)}/`),
  },
);
