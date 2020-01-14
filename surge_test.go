package surge_test

import (
	"bytes"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/renproject/surge"
)

var _ = Describe("Surge", func() {
	Context("when marshaling int8s", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				x := int8(rand.Intn(256))

				buf := new(bytes.Buffer)
				Expect(Marshal(x, buf)).ToNot(HaveOccurred())

				y := int8(0)
				Expect(Unmarshal(&y, buf)).ToNot(HaveOccurred())

				Expect(x).To(Equal(y))
			})
		})
	})

	Context("when marshaling *int8 pointers", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				x := int8(rand.Intn(256))

				buf := new(bytes.Buffer)
				Expect(Marshal(&x, buf)).ToNot(HaveOccurred())

				y := int8(0)
				Expect(Unmarshal(&y, buf)).ToNot(HaveOccurred())

				Expect(x).To(Equal(y))
			})
		})
	})

	Context("when marshaling []int8 arrays", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				xs := [1000000]int8{}
				for i := 0; i < len(xs); i++ {
					xs[i] = int8(rand.Intn(256))
				}

				buf := new(bytes.Buffer)
				Expect(Marshal(xs, buf)).ToNot(HaveOccurred())

				ys := []int8{}
				Expect(Unmarshal(&ys, buf)).ToNot(HaveOccurred())

				Expect(xs[:]).To(Equal(ys))
			})
		})
	})

	Context("when marshaling *[]int8 array pointers", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				xs := [1000000]int8{}
				for i := 0; i < len(xs); i++ {
					xs[i] = int8(rand.Intn(256))
				}

				buf := new(bytes.Buffer)
				Expect(Marshal(&xs, buf)).ToNot(HaveOccurred())

				ys := []int8{}
				Expect(Unmarshal(&ys, buf)).ToNot(HaveOccurred())

				Expect(xs[:]).To(Equal(ys))
			})
		})
	})

	Context("when marshaling []int8 slices", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				xs := make([]int8, rand.Intn(1000000))
				for i := 0; i < len(xs); i++ {
					xs[i] = int8(rand.Intn(256))
				}

				buf := new(bytes.Buffer)
				Expect(Marshal(xs, buf)).ToNot(HaveOccurred())

				ys := []int8{}
				Expect(Unmarshal(&ys, buf)).ToNot(HaveOccurred())

				Expect(xs).To(Equal(ys))
			})
		})
	})

	Context("when marshaling *[]int8 slice pointers", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				xs := make([]int8, rand.Intn(1000000))
				for i := 0; i < len(xs); i++ {
					xs[i] = int8(rand.Intn(256))
				}

				buf := new(bytes.Buffer)
				Expect(Marshal(&xs, buf)).ToNot(HaveOccurred())

				ys := []int8{}
				Expect(Unmarshal(&ys, buf)).ToNot(HaveOccurred())

				Expect(xs).To(Equal(ys))
			})
		})
	})

	Context("when marshaling map[int8]int maps", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				xs := map[int8]int8{}
				len := rand.Intn(1000000)
				for i := 0; i < len; i++ {
					xs[int8(rand.Intn(256))] = int8(rand.Intn(256))
				}

				buf := new(bytes.Buffer)
				Expect(Marshal(xs, buf)).ToNot(HaveOccurred())

				ys := map[int8]int8{}
				Expect(Unmarshal(&ys, buf)).ToNot(HaveOccurred())

				Expect(xs).To(Equal(ys))
			})
		})
	})

	Context("when marshaling structs", func() {
		Context("when marshaling and then unmarshaling", func() {
			It("should equal itself", func() {
				type Person struct {
					Name       string   `surge:"0"`
					Age        uint64   `surge:"1"`
					Friends    []string `surge:"2"`
					PubIgnore  bool
					privIgnore bool
				}
				x := Person{
					Name:       "Alice",
					Age:        25,
					Friends:    []string{"Bob", "Eve"},
					PubIgnore:  true,
					privIgnore: true,
				}

				buf := new(bytes.Buffer)
				Expect(Marshal(x, buf)).ToNot(HaveOccurred())

				y := Person{}
				Expect(Unmarshal(&y, buf)).ToNot(HaveOccurred())
				Expect(y.PubIgnore).To(BeFalse())
				Expect(y.privIgnore).To(BeFalse())
				y.PubIgnore = true
				y.privIgnore = true

				Expect(x).To(Equal(y))
			})

			Context("when structs are recursive", func() {
				It("should equal itself", func() {
					type Person struct {
						Name    string   `surge:"0"`
						Age     uint64   `surge:"1"`
						Friends []Person `surge:"2"`
					}

					x := Person{
						Name: "Alice",
						Age:  25,
						Friends: []Person{
							Person{
								Name:    "Bob",
								Age:     26,
								Friends: []Person{},
							},
						},
					}

					buf := new(bytes.Buffer)
					Expect(Marshal(x, buf)).ToNot(HaveOccurred())

					y := Person{}
					Expect(Unmarshal(&y, buf)).ToNot(HaveOccurred())

					Expect(x).To(Equal(y))
				})
			})
		})

	})
})
