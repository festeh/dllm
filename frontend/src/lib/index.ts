// place files you want to import through the `$lib` alias in this folder.
//

export function resizeOutputField(target: HTMLElement) {
  if (target.clientHeight >= target.scrollHeight && target.clientHeight < 50) {
    return;
  }
  target.style.height = '10px';
  target.style.height = +target.scrollHeight + 'px';
}
