// casper.test.begin('Google search retrieves 10 or more results', 5, function suite(test) {
//     casper.start("http://www.google.fr/", function() {
//         test.assertTitle("Google", "google homepage title is the one expected");
//         test.assertExists('form[action="/search"]', "main form is found");
//         this.fill('form[action="/search"]', {
//             q: "casperjs"
//         }, true);
//     });

//     casper.then(function() {
//         test.assertTitle("casperjs - Recherche Google", "google title is ok");
//         test.assertUrlMatch(/q=casperjs/, "search term has been submitted");
//         test.assertEval(function() {
//             return __utils__.findAll("h3.r").length >= 10;
//         }, "google search for \"casperjs\" retrieves 10 or more results");
//     });

//     casper.run(function() {
//         test.done();
//     });
// });

casper.test.begin('Test creating new link post', function suite(test) {
    casper.start('http://localhost:3000/posts/new', function() {
        loginIfNeeded(function() {
            test.assertExists('#submit-form-tab-link form', "new post form is found");
            casper.fill('#submit-form-tab-link form', {
                title: 'test title',
                link: 'http://example.com',
            }, true);

            casper.then(function() {
                test.assertUrlMatch(/\/post\/view\/(.*)/, "post has been created");
                test.assertSelectorHasText(".post-heading", "test title")
            });
        });
    })

    casper.run(function() {
        test.done();
    });
});

casper.test.begin('Test creating new link post with no title', function suite(test) {
    casper.start('http://localhost:3000/posts/new', function() {
        loginIfNeeded(function() {
            test.assertExists('#submit-form-tab-link form', "new post form is found");
            casper.fill('#submit-form-tab-link form', {
                title: '',
                link: 'http://example.com',
            }, true);

            casper.then(function() {
                test.assertSelectorHasText(".help-block", "Required")
            });
        });
    })

    casper.run(function() {
        test.done();
    });
});

function loginIfNeeded(cb) {
    if (casper.exists('form[action="/login"]')) {
        casper.echo('authenticating as test user')
        casper.fill('form[action="/login"]', {
            username: 'test',
            password: 'password',
        }, true);

        casper.then(function() {
            cb();
        });
    } else {
        cb();
    }
}
