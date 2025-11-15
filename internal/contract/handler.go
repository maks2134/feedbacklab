package contract

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ContractHandler struct {
	service ContractService
}

func NewContractHandler(service *ContractService) *ContractHandler {
	return &ContractHandler{service: *service}
}

// GetAll godoc
// @Summary Get all contracts
// @Description Returns all contracts
// @Tags Contracts
// @Produce json
// @Success 200 {array} models.Contract
// @Router /api/contracts [get]
func (h *ContractHandler) GetAll(c *fiber.Ctx) error {
	items, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(items)
}

// GetByID godoc
// @Summary Get contract by ID
// @Description Returns a single contract by its ID
// @Tags Contracts
// @Produce json
// @Param id path int true "Contract ID"
// @Success 200 {object} models.Contract
// @Router /api/contracts/{id} [get]
func (h *ContractHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	item, err := h.service.GetByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found"})
	}
	return c.JSON(item)
}

// Create godoc
// @Summary Create a new contract
// @Tags Contracts
// @Accept json
// @Produce json
// @Param contract body models.Contract true "Contract Data"
// @Success 201 {object} models.Contract
// @Router /api/contracts [post]
func (h *ContractHandler) Create(c *fiber.Ctx) error {
	var input Contract
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.service.Create(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(input)
}

// Update godoc
// @Summary Update existing contract
// @Tags Contracts
// @Accept json
// @Produce json
// @Param id path int true "Contract ID"
// @Param contract body models.Contract true "Updated Contract Data"
// @Success 200 {object} models.Contract
// @Router /api/contracts/{id} [put]
func (h *ContractHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	var input Contract
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	input.ID = id
	if err := h.service.Update(&input); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(input)
}

// Delete godoc
// @Summary Delete contract
// @Tags Contracts
// @Param id path int true "Contract ID"
// @Success 204 "No Content"
// @Router /api/contracts/{id} [delete]
func (h *ContractHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.service.Delete(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(204)
}
