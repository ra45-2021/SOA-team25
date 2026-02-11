using TourService.Contracts;
using TourService.Domain;
using TourService.Repositories;

namespace TourService.Services;

public class TourService : ITourService
{
    private readonly ITourRepository _tourRepo;
    private readonly ITagRepository _tagRepo;

    public TourService(ITourRepository tourRepo, ITagRepository tagRepo)
    {
        _tourRepo = tourRepo;
        _tagRepo = tagRepo;
    }

    public async Task<TourDto> CreateAsync(string authorId, CreateTourRequest req, CancellationToken ct)
    {
        if (string.IsNullOrWhiteSpace(req.Name)) throw new ArgumentException("Name is required.");
        if (string.IsNullOrWhiteSpace(req.Description)) throw new ArgumentException("Description is required.");

        var tags = await _tagRepo.GetOrCreateAsync(req.Tags, ct);

        var now = DateTime.UtcNow;

        var tour = new Tour
        {
            Name = req.Name.Trim(),
            Description = req.Description.Trim(),
            Difficulty = req.Difficulty,
            Status = TourStatus.Draft,
            Price = 0m,
            AuthorId = authorId,
            CreatedAtUtc = now,
            UpdatedAtUtc = now,
            TourTags = tags.Select(t => new TourTag { Tag = t }).ToList()
        };

        var created = await _tourRepo.AddAsync(tour, ct);
        return Map(created);
    }

    public async Task<List<TourDto>> GetMyToursAsync(string authorId, CancellationToken ct)
    {
        var tours = await _tourRepo.GetByAuthorAsync(authorId, ct);
        return tours.Select(Map).ToList();
    }

    private static TourDto Map(Tour t)
    {
        return new TourDto
        {
            Id = t.Id,
            Name = t.Name,
            Description = t.Description,
            Difficulty = t.Difficulty,
            Status = t.Status,
            Price = t.Price,
            Tags = t.TourTags.Select(x => x.Tag.Name).OrderBy(x => x).ToList(),
            CreatedAtUtc = t.CreatedAtUtc,
            UpdatedAtUtc = t.UpdatedAtUtc
        };
    }
}
