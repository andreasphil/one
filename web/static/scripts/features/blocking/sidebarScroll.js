(function init(id = "sidebar") {
  const sidebar = document.getElementById(id);
  const key = `${id}-scroll`;

  if (!(sidebar instanceof HTMLElement)) {
    console.warn(`could not find sidebar element with ID ${id}, scroll position not restored`);
    return;
  }

  if (sidebar.scrollHeight > sidebar.clientHeight) {
    const saved = sessionStorage.getItem(key);

    if (saved !== null) {
      sidebar.scrollTop = Number(saved);
    } else {
      sidebar.querySelector("[aria-current]")?.scrollIntoView({ block: "center" });
    }

    sidebar.addEventListener("scrollend", () => {
      sessionStorage.setItem(key, sidebar.scrollTop.toString());
    });
  }
})();
