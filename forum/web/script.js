function likePost(button) {
    let likeCount = button.querySelector('.like-count');
    likeCount.textContent = parseInt(likeCount.textContent) + 1;
}

function dislikePost(button) {
    let dislikeCount = button.querySelector('.dislike-count');
    dislikeCount.textContent = parseInt(dislikeCount.textContent) + 1;
}

function toggleComments(button) {
    let commentsSection = button.parentElement.nextElementSibling;
    if (commentsSection.style.display === 'none') {
        commentsSection.style.display = 'block';
    } else {
        commentsSection.style.display = 'none';
    }
}

function addComment(button) {
    let commentInput = button.previousElementSibling.previousElementSibling;
    let commentText = button.previousElementSibling;
    let commentAuthor = commentInput.value;
    let commentContent = commentText.value;

    if (commentAuthor && commentContent) {
        let commentsList = button.parentElement.querySelector('.comments-list');
        let newComment = document.createElement('li');

        let authorSpan = document.createElement('span');
        authorSpan.className = 'comment-author';
        authorSpan.textContent = commentAuthor;

        let contentSpan = document.createElement('span');
        contentSpan.textContent = commentContent;

        let metaSpan = document.createElement('span');
        metaSpan.className = 'comment-meta';
        metaSpan.textContent = `on ${new Date().toLocaleString()}`;

        newComment.appendChild(authorSpan);
        newComment.appendChild(contentSpan);
        newComment.appendChild(metaSpan);

        commentsList.appendChild(newComment);

        commentInput.value = '';
        commentText.value = '';
    }
}

function filterPosts() {
    let filter = document.getElementById('category-filter').value;
    let posts = document.querySelectorAll('.post');
    posts.forEach(post => {
        if (filter === 'all' || post.getAttribute('data-category') === filter) {

            post.style.display = 'block';
        } else {
            post.style.display = 'none';
        }
    });
}

function navigate(page) {
    window.location.href = page;
}
