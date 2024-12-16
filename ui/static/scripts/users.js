// Automatically submit form on select value changes
const selectUser = document.getElementById('select-user')
const selectUserForm = document.getElementById('select-user-form')
if (selectUser && selectUserForm) {
    selectUser.addEventListener('change', () => selectUserForm.submit())
}

// Require delete form confirmation
const deleteUserForm = document.getElementById('delete-user-form')
if (deleteUserForm) {
    deleteUserForm.addEventListener('submit', (e) => {
        if (!confirm("Confirm user deletion. This can not be undone.")) {
            e.preventDefault()
        }
    })
}
