"use strict";
var KTSigninGeneral = function () {
    var e, t, i;
    return {
        init: function () {
            e = document.querySelector("#kt_sign_in_form"), t = document.querySelector("#kt_sign_in_submit"), i = FormValidation.formValidation(e, {
                fields: {
                    email: {
                        validators: {
                            notEmpty: {
                                message: "Email tidak boleh kosong"
                            },
                            // regexp: {
                            //     regexp: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                            //     message: "Email tidak valid"
                            // }
                        }
                    },
                    password: {
                        validators: {
                            notEmpty: {
                                message: "Password tidak boleh kosong"
                            },
                            // stringLength: {
                            //     min: 8,
                            //     message: 'Password minimal 8 karakter'
                            // }
                        }
                    }
                },
                plugins: {
                    trigger: new FormValidation.plugins.Trigger(),
                    bootstrap: new FormValidation.plugins.Bootstrap5({
                        rowSelector: ".fv-row",
                        eleInvalidClass: "",
                        eleValidClass: ""
                    })
                }
            }), t.addEventListener("click", (function (n) {
                n.preventDefault();
                var form = document.querySelector("#kt_sign_in_form");
                i.validate().then((function (i) {
                    if ("Valid" == i) {
                        t.setAttribute("data-kt-indicator", "on");
                        t.disabled = !0;
                        var xhr = new XMLHttpRequest();
                        var url = e.getAttribute("action");
                        var form_data = new FormData(e);
                        xhr.open(e.getAttribute('method'), url);
                        xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
                        xhr.send(new URLSearchParams([...form_data.entries()]).toString());
                        xhr.onload = function () {
                            if (xhr.status === 200) {
                                // t.removeAttribute("data-kt-indicator");
                                // t.disabled = !1;
                                // Swal.fire({
                                //     text: "Berhasil login",
                                //     icon: "success",
                                //     buttonsStyling: !1,
                                //     confirmButtonText: "Ok",
                                //     customClass: {
                                //         confirmButton: "btn btn-primary"
                                //     }
                                // })
                                //console.log(response.clientid + "?access_token=" + response.access_token + "&user=" + response.user)
                        
                            window.location.replace(window.location.origin + "/dashboard");
                            
                               
                            } else {
                                console.log(xhr)
                                const response = JSON.parse(xhr.responseText);
                                if (xhr.status === 204) {
                                    window.location.replace("/")
                                }
                                t.removeAttribute("data-kt-indicator");
                                t.disabled = !1;
                                Swal.fire({
                                    title: response.error_title,
                                    text: response.error_text,
                                    icon: "error",
                                    buttonsStyling: !1,
                                    confirmButtonText: "Ok, mengerti",
                                    customClass: {
                                        confirmButton: "btn btn-primary"
                                    }
                                });
                            }
                        };
                        xhr.onerror = function () {
                            t.removeAttribute("data-kt-indicator");
                            t.disabled = !1;
                            Swal.fire({
                                text: "Sorry, looks like there are some errors detected, please try again.",
                                icon: "error",
                                buttonsStyling: !1,
                                confirmButtonText: "Ok, paham",
                                customClass: {
                                    confirmButton: "btn btn-primary"
                                }
                            });
                        };
                    }
                    //  else {
                    //     Swal.fire({
                    //         text: "Masih ada field yang kosong. Pastikan semua field terisi.",
                    //         icon: "error",
                    //         buttonsStyling: !1,
                    //         confirmButtonText: "Ok",
                    //         customClass: {
                    //             confirmButton: "btn btn-primary"
                    //         }
                    //     });
                    // }
                }))
            }))
        }
    }
}();
KTUtil.onDOMContentLoaded((function () {
    KTSigninGeneral.init()
}));