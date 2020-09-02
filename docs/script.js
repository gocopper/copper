var anchors = document.getElementsByClassName('anchor');
var links = document.getElementsByTagName('a');

setActiveLink();
window.addEventListener('scroll', debounce(setURLHash, 300), false);

function setURLHash() {
  var a = currentAnchor();

  history.pushState(null, null, '#' + a.getAttribute('id'));
  setActiveLink();
}

function setActiveLink() {
  var defaultHash = '#introduction';

  for (var i=0; i<links.length; i++) {
    var link = links[i];

    link.classList.remove('active');

    if (link.getAttribute('href') === (location.hash || defaultHash)) {
      link.classList.add('active');
    }
  }
}

function currentAnchor() {
  var pageY = window.pageYOffset;

  for (var i=0; i<anchors.length-1; i++) {
    var y = pageY + 300;
    var currAnchorY = anchors[i].getBoundingClientRect().y + pageY;
    var nextAnchorY = anchors[i+1].getBoundingClientRect().y + pageY;

    if (y >= currAnchorY && y < nextAnchorY) {
      return anchors[i];
    }
  }

  return pageY <= 100 ? anchors[0] : anchors[anchors.length-1];
}

// https://gist.github.com/peduarte/7ee475dd0fae1940f857582ecbb9dc5f#file-index-js
function debounce(func, wait) {
  let timeout;
  return function(...args) {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      func.apply(this, args);
    }, wait);
  };
}
