OUTPUT := a.out

# compiler options
CC := gcc
LINKER := gcc
INCDIRS := -I.
CFLAGS := -Wall #-Wextra

# generate list of source, object, and dependency files
SRCFILES := $(wildcard *.c)
OBJFILES := $(patsubst %.c,%.o,$(SRCFILES))
DEPFILES := $(patsubst %.c,%.d,$(SRCFILES))

# rule to build output
$(OUTPUT): $(OBJFILES)
	    $(LINKER) $^ -o $@

# compile object files
%.o: %.c
	    $(CC) $(CFLAGS) $(INCDIRS) -c $< -o $@

# generate dependency files
%.d: %.c
	    $(CC) $(INCDIRS) -MM $< > $@

# include the dependency files
-include $(DEPFILES)
