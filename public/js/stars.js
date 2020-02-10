const ANIMATION_STAR_SPEED = 500;

function star(rating, tag) {

  var t = document.querySelector(tag)
  var div = document.createElement("div")
  for (i = 0;i < rating;i++) {
    setTimeout(function (i) {
      let INNER = t.querySelector("div")
      if (i < rating) {
        let el = document.createElement("span")
        el.setAttribute('class', 'uk-icon ml2 rating-bar uk-animation-slide-left')
        el.setAttribute('uk-icon', 'star')
        INNER.appendChild(el);
      }
      else return
    }, i * ANIMATION_STAR_SPEED, i)
  }
  t.prepend(div)
}


document.addEventListener('DOMContentLoaded', function () {
  let reviews = document.querySelectorAll('.uk-card')
  reviews.forEach(review => {
    let r = review.querySelector('.uk-card-body').querySelector('.uk-text-meta > a')
    if (r !== null) {
      r.innerText = r.innerText.replace(/-/g, " ").toUpperCase()
    }

  })
  let stars = document.querySelectorAll('.star.starRating')
  stars.forEach(star => {
    star.checked = true
  })
  let currentRating = 5;
  stars.forEach(star => {
    star.addEventListener('click', (e) => {
      currentRating = e.target.value
      if (currentRating === '5') {
        stars[4].checked = true
        stars[3].checked = true
        stars[2].checked = true
        stars[1].checked = true
        stars[0].checked = true
      } else if (currentRating === '4') {
        // remove check on higher values
        stars[4].checked = false
        // add current values
        stars[3].checked = true
        stars[2].checked = true
        stars[1].checked = true
        stars[0].checked = true
      } else if (currentRating === '3') {
        // remove check on higher values
        stars[4].checked = false
        stars[3].checked = false
        // add current values
        stars[2].checked = true
        stars[1].checked = true
        stars[0].checked = true
      } else if (currentRating === '2') {
        // remove check on higher values
        stars[4].checked = false
        stars[3].checked = false
        stars[2].checked = false
        // add current values
        stars[1].checked = true
        stars[0].checked = true
      } else if (currentRating === '1') {
        // remove check on higher values
        stars[4].checked = false
        stars[3].checked = false
        stars[2].checked = false
        stars[1].checked = false
        // add current values
        stars[0].checked = true
      }
    })
  })
})
