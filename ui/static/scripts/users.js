// Automatically submit select user form on select value changes
const selectUser = document.getElementById('select-user')
const selectUserForm = document.getElementById('select-user-form')
if (selectUser && selectUserForm) {
    selectUser.addEventListener('change', () => selectUserForm.submit())
}

// Automatically submit permissions form on input value changes
const permissionsForm = document.getElementById('permissions-form')
if (permissionsForm) {
    const permissionsInputs = permissionsForm.querySelectorAll('input')
    permissionsInputs.forEach(input => {
        input.addEventListener('change', () => permissionsForm.submit())
    })
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
