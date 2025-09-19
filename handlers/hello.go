package handlers

import "snaczek-server/coreutils"

func HelloHandler(req coreutils.Request) coreutils.Response {
	return coreutils.Response{
		Status_code: 200,
		Headers: map[string]string{"Content-Type": "text/plain"},
		Body: []byte("Hello from GET /hello\n"),
	}
}

func JsonTestHandler(req coreutils.Request) coreutils.Response {
    body := []byte(`
{
  "name": "Mateusz",
  "surname": "Filoda",
  "email": "mateusz.filoda@example.com",
  "birth_date": "2002-05-14",
  "created_at": "2025-09-04T12:34:56Z",
  "description": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum in gravida magna, et maximus lorem. Curabitur sit amet arcu nec risus sagittis laoreet. Donec interdum justo a elit pretium, a vulputate dui elementum. Sed accumsan libero nec risus blandit, eget scelerisque felis efficitur. Pellentesque at sapien sit amet nunc vestibulum faucibus. Nam pharetra, velit nec suscipit laoreet, dolor neque aliquam quam, sit amet tincidunt est libero at nibh. Integer ullamcorper sapien nec augue efficitur, nec iaculis nibh pellentesque. Aliquam erat volutpat. Duis malesuada tellus vel malesuada ultricies. Nulla facilisi. Quisque ut nunc vitae quam egestas feugiat. Sed nec eros erat. Integer ullamcorper, turpis vitae dignissim malesuada, nunc ligula cursus eros, sed congue massa nulla ac turpis. Praesent ullamcorper, ipsum sit amet rhoncus faucibus, nibh arcu consequat sapien, nec convallis neque mi ac libero. Vestibulum et velit varius, tristique ligula et, aliquam turpis. Cras rutrum massa sit amet ante euismod, eget bibendum nisl fermentum.",
  "key": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6Ik1hdGV1c3oiLCJpYXQiOjE1MTYyMzkwMjJ9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
  "friends": [
    {
      "name": "Anna",
      "surname": "Kowalska",
      "email": "anna.k@example.com",
      "birth_date": "2003-02-11"
    },
    {
      "name": "Piotr",
      "surname": "Nowak",
      "email": "piotr.nowak@example.com",
      "birth_date": "2001-11-23"
    },
    {
      "name": "Kasia",
      "surname": "Mazur",
      "email": "kasia.mazur@example.com",
      "birth_date": "2002-08-19"
    },
    {
      "name": "Tomasz",
      "surname": "Lewandowski",
      "email": "t.lewandowski@example.com",
      "birth_date": "2000-07-04"
    },
    {
      "name": "Ola",
      "surname": "Wójcik",
      "email": "ola.wojcik@example.com",
      "birth_date": "2004-01-29"
    },
    {
      "name": "Bartek",
      "surname": "Krawczyk",
      "email": "bartek.krawczyk@example.com",
      "birth_date": "1999-12-09"
    },
    {
      "name": "Zuzanna",
      "surname": "Kamińska",
      "email": "zuzanna.kam@example.com",
      "birth_date": "2003-06-15"
    },
    {
      "name": "Marek",
      "surname": "Dąbrowski",
      "email": "marek.dab@example.com",
      "birth_date": "2002-03-21"
    },
    {
      "name": "Ewa",
      "surname": "Pawlak",
      "email": "ewa.pawlak@example.com",
      "birth_date": "2001-10-05"
    },
    {
      "name": "Jan",
      "surname": "Zieliński",
      "email": "jan.ziel@example.com",
      "birth_date": "2000-09-12"
    }
  ]
}
	`)

    return coreutils.Response{
        Status_code: 200,
        Headers: map[string]string{"Content-Type": "application/json"},
        Body: body,
    }
}
