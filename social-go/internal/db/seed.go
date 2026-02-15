package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"social/internal/store"
)

var userNames = []string{
	"Ana", "Bruno", "Carla", "Daniel", "Elena", "Federico", "Gabriela", "Hernan", "Ines", "Julian",
	"Karina", "Luciano", "Maria", "Nicolas", "Olivia", "Pablo", "Queralt", "Ramiro", "Sofia", "Tomas",
	"Ulises", "Valentina", "Walter", "Ximena", "Yamila", "Zoe", "Andres", "Belen", "Camila", "Damian",
	"Emilia", "Facundo", "Guadalupe", "Horacio", "Ignacio", "Jimena", "Kevin", "Lautaro", "Martina", "Nahuel",
	"Octavio", "Paula", "Renata", "Santiago", "Thiago", "Ursula", "Victoria", "Wanda", "Xavier", "Yesica",

	"Agustin", "Barbara", "Cecilia", "Diego", "Esteban", "Florencia", "Gonzalo", "Helena", "Ivana", "Joaquin",
	"Kiara", "Leandro", "Milagros", "Norberto", "Oriana", "Patricio", "Rocio", "Sebastian", "Tatiana", "Valerio",
	"Wilmer", "Yasmin", "Zaira", "Alberto", "Beatriz", "Claudio", "Delfina", "Ezequiel", "Fiorella", "Gerardo",
	"Hugo", "Iara", "Jose", "Lidia", "Manuel", "Noelia", "Omar", "Priscila", "Rafael", "Sol",
	"Teo", "Violeta", "Wendy", "Axel", "Bianca", "Ciro", "Dario", "Ema", "Franco", "Gisela",
}

var titles = []string{
	"Getting Started with Go",
	"Understanding Goroutines",
	"Mastering Channels in Go",
	"Building REST APIs in Go",
	"Dependency Injection in Go",
	"Error Handling Best Practices",
	"Structs and Interfaces Explained",
	"Writing Clean Go Code",
	"Testing in Go Made Simple",
	"Concurrency Patterns in Go",
	"Working with JSON in Go",
	"Middleware in Go Web Apps",
	"Context Usage in Go",
	"Building CLI Tools in Go",
	"Optimizing Go Performance",
	"Logging Strategies in Go",
	"Graceful Shutdown in Go",
	"Microservices with Go",
	"Using Generics in Go",
	"Go Project Structure Tips",
}

var contents = []string{
	"Go is a statically typed, compiled language designed for simplicity and performance. Its concurrency model makes it ideal for scalable systems.",
	"In this article, we explore how goroutines allow lightweight concurrency and how channels help coordinate communication safely.",
	"Building REST APIs in Go is straightforward using net/http or popular routers like chi and gorilla/mux.",
	"Understanding interfaces in Go is key to writing flexible and testable code.",
	"Error handling in Go is explicit by design, encouraging developers to handle failures properly.",
	"Using context.Context correctly prevents memory leaks and improves cancellation handling in distributed systems.",
	"Struct composition is preferred over inheritance in Go, leading to cleaner and more maintainable code.",
	"Testing in Go is built-in, with the testing package offering powerful tools for unit and integration tests.",
	"JSON marshaling and unmarshaling is simple with the encoding/json package.",
	"Dependency injection in Go is often done manually, promoting clarity over magic frameworks.",
	"Graceful shutdown in Go servers ensures no requests are lost during deployment.",
	"Generics introduced in Go 1.18 allow writing reusable and type-safe functions.",
	"Channels can be buffered or unbuffered, depending on synchronization requirements.",
	"Go modules simplify dependency management and version control.",
	"Logging should be structured and consistent, especially in microservices architectures.",
	"Benchmarking Go applications helps identify performance bottlenecks early.",
	"Middleware patterns enable reusable logic in HTTP applications.",
	"Concurrency bugs can be detected using the race detector with the -race flag.",
	"CLI applications in Go are easy to build using libraries like cobra.",
	"Go’s simplicity makes it a strong choice for cloud-native development.",
}

var tags = []string{
	"go",
	"golang",
	"backend",
	"api",
	"rest",
	"microservices",
	"concurrency",
	"goroutines",
	"channels",
	"interfaces",
	"testing",
	"performance",
	"cloud",
	"docker",
	"kubernetes",
	"devops",
	"clean-code",
	"architecture",
	"json",
	"generics",
}

var comments = []string{
	"Great article! The explanation about goroutines was very clear.",
	"I’ve been struggling with channels, this helped a lot.",
	"Could you write more about real-world concurrency patterns?",
	"This clarified many doubts I had about interfaces.",
	"Very practical examples, thanks for sharing.",
	"I’d love to see a follow-up post about performance tuning.",
	"The section about context usage was especially helpful.",
	"Nice breakdown of error handling best practices.",
	"Go keeps surprising me with its simplicity.",
	"Do you recommend any advanced Go resources?",
	"I tried this approach in my project and it worked perfectly.",
	"Clear, concise, and easy to understand.",
	"More posts like this please!",
	"This saved me hours of debugging.",
	"Can you compare this with similar patterns in Java?",
	"I appreciate the clean examples included.",
	"Helpful overview for beginners.",
	"I didn’t know about the race detector flag, great tip!",
	"Looking forward to more content on microservices.",
	"Thanks for making complex topics approachable.",
}

func Seed(store store.Storage) error {
	ctx := context.Background()

	users := generateUsers(100)

	for i, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Printf("Error creating user at i=%d username=%q email=%q err=%v", i, user.Username, user.Email, err)
			return err
		}
	}

	post := generatePosts(200, users)

	for _, post := range post {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating post", err)
			return nil
		}
	}

	comments := generateComments(500, users, post)

	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comment:", err)
			return nil
		}
	}

	log.Println("Seeding complete")

	return nil
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: userNames[i] + fmt.Sprintf("%d", i),
			Email:    userNames[i] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)
	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)
	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
