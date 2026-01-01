(() => {
  const paint = performance
    .getEntriesByType("paint")
    .find((p) => p.name === "first-contentful-paint");
  const fcp = paint ? paint.startTime : 0;

  const origin = window.location.origin;

  const allUrls = Array.from(document.querySelectorAll("a[href]")).map(
    (a) => a.href,
  );

  const internalUrls = allUrls.filter((href) => href.startsWith(origin));

  // remove duplicates
  const uniqueLinks = internalUrls.filter((v, i, a) => a.indexOf(v) === i);

  return JSON.stringify({
    fcp: fcp,
    links: uniqueLinks,
  });
})();
