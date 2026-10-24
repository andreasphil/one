(function init(id = "navigation") {
  const navigation = document.getElementById(id);
  const key = `${id}-scroll`;

  if (!(navigation instanceof HTMLElement)) {
    console.warn(`could not find navigation element with ID ${id}, scroll position not restored`);
    return;
  }

  if (navigation.scrollHeight > navigation.clientHeight) {
    const saved = sessionStorage.getItem(key);

    if (saved !== null) {
      navigation.scrollTop = Number(saved);
    } else {
      navigation.querySelector("[aria-current]")?.scrollIntoView({ block: "center" });
    }

    navigation.addEventListener("scrollend", () => {
      sessionStorage.setItem(key, navigation.scrollTop.toString());
    });
  }
})();
