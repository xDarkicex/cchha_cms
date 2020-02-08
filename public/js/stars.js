
function star(rating, tag) {
  let review = document.querySelector(tag)
  console.log(tag)
  console.log(review)

  let ratingvalue = rating;
  let t = review.querySelector('div.p1 > div.uk-card')
  // let thing = document.querySelector('.uk-text-meta > a')
  // console.log(thing)
  // console.log(t)

  // review.querySelector('.uk-text-meta > a').innerText.replace(/-/g, " ").toUpperCase()


  var div = document.createElement('div')

  for (i = 0;i < rating;i++) {
    var el = document.createElement('span')
    el.setAttribute('class', 'uk-icon ml2 rating-bar')
    el.setAttribute('uk-icon', 'star')
    div.appendChild(el)
  }
  t.prepend(div)
}


document.addEventListener('DOMContentLoaded', function () {
  let reviews = document.querySelectorAll('.uk-card')
  reviews.forEach(review => {
    let rating_data = review.querySelector('p.rating')
    var r = rating_data.innerText
    rating_data.parentNode.removeChild(rating_data)
    review.querySelector('.uk-text-meta > a').innerText = review.querySelector('.uk-text-meta > a').innerText.replace(/-/g, " ").toUpperCase()
    for (i = 0;i < r;i++) {
      var el = document.createElement('span')
      el.setAttribute('class', 'uk-icon ml2 rating-bar')
      el.setAttribute('uk-icon', 'star')
      review.prepend(el)
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
