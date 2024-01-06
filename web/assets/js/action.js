function togglePassword(id) {
    const passwordInput = document.getElementById(id);
    const passwordToggle = document.getElementById('password-toggle');

    if (passwordInput.type === "password") {
        passwordInput.type = "text";
        passwordToggle.classList.remove('bi-eye');
        passwordToggle.classList.add('bi-eye-slash');
    } else {
        passwordInput.type = "password";
        passwordToggle.classList.remove('bi-eye-slash');
        passwordToggle.classList.add('bi-eye');
    }
}

// Stepper lement
// var element = document.querySelector("#kt_stepper_example_basic");

// // Initialize Stepper
// var stepper = new KTStepper(element);

// // Handle next step
// stepper.on("kt.stepper.next", function (stepper) {
//     stepper.goNext(); // go next step
// });

// // Handle previous step
// stepper.on("kt.stepper.previous", function (stepper) {
//     stepper.goPrevious(); // go previous step
// });