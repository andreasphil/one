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

const notes = JSON.parse(document.getElementById("clientState.NotesMeta").textContent).map((i) => ({
  id: i.Slug,
  name: i.Title,
  groupName: "Notes",
  icon: icons.get("StickyNote").cloneNode(true),
  action: () => {
    navigation.navigate(`/notes/${i.Slug}/`);
  },
}));

CommandBar.define();
CommandBar.instance.registerCommand(...notes);
