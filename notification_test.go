package ant

import "testing"

func TestChat(t *testing.T) {
	c := &Chat{
		URL: "https://chat.googleapis.com/v1/spaces/AAAAYFDQezs/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=Zevs_4leRsG8C8kmBqMUNZH9ioP3ES9yzkPWLM0-mz8%3D",
	}
	if err := c.SendText(map[string]string{
		"hello bot": "aaa",
	}); err != nil {
		t.Fatal(err)
	}
	if err := c.SendCard("https://txthinking-file.storage.googleapis.com/feada7287db64e3d9fd838fb38c13227%2F%E6%9C%80%E6%96%B0%E6%9C%AC%E5%9C%9F%E5%A4%A7%E5%A5%B6%E5%A6%B9%E5%81%9A%E6%84%9B%E7%9F%AD%E7%89%87%EF%BC%81%E8%BA%AB%E6%9D%90%E8%B6%85%E6%AD%A3.jpeg", "最新本土大奶妹做愛短片！身材超正", "https://txthinking-file.storage.googleapis.com/4e7087f8c73d4deda0dd7bc12c040d3f/%E6%9C%80%E6%96%B0%E6%9C%AC%E5%9C%9F%E5%A4%A7%E5%A5%B6%E5%A6%B9%E5%81%9A%E6%84%9B%E7%9F%AD%E7%89%87%EF%BC%81%E8%BA%AB%E6%9D%90%E8%B6%85%E6%AD%A3.mp4"); err != nil {
		t.Fatal(err)
	}
}
